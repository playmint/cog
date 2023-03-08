import Phaser from 'phaser';
import { ApolloClient, NormalizedCacheObject, InMemoryCache, gql, split, HttpLink, FetchResult } from '@apollo/client';
import { getMainDefinition } from '@apollo/client/utilities';
import { createClient } from "graphql-ws";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { BigNumber } from "ethers";
import * as ethers from "ethers";

const COG_WS_ENDPOINT = import.meta.env.VITE_COG_WS_ENDPOINT || 'ws://localhost:8080/query';
const COG_HTTP_ENDPOINT = import.meta.env.VITE_COG_HTTP_ENDPOINT || 'http://localhost:8080/query';

const STATE_FRAGMENT = `
    block
    tiles: nodes(match: {kinds: ["Tile"]}) {
        coords: keys
        biome: value(match: {via: [{rel: "Biome"}]})
        seed: node(match: {kinds: ["Seed"], via: [{rel: "ProvidesEntropyTo", dir: IN}]}) {
          key
        }
    }
    seekers: nodes(match: {kinds: ["Seeker"]}) {
        key
        position: node(match: {kinds: ["Tile"], via:[{rel: "Location"}]}) {
            coords: keys
        }
        player: node(match: {kinds: ["Player"], via:[{rel: "Owner"}]}) {
            address: key
        }
        cornBalance: value(match: {via: [{rel: "Balance"}]})
    }
`;

const STATE_SUBSCRIPTION = gql(`
    subscription OnState {
        state(gameID: "latest") {
            ${STATE_FRAGMENT}
        }
    }
`);

const TX_SUBSCRIPTION = gql(`
    subscription OnTx($owner: String!) {
        transaction(gameID: "latest", owner: $owner) {
            id
            status
        }
    }
`);

const STATE_QUERY = gql(`
    query GetState {
        game(id: "latest") {
            state {
                ${STATE_FRAGMENT}
            }
        }
    }
`);

const DISPATCH = gql(`
    mutation dispatch($gameID: ID!, $actions: [String!]!, $auth: String!) {
        dispatch(
            gameID: $gameID,
            actions: $actions,      # encoded action bytes
            authorization: $auth    # session's signature of $action
        ) {
            id
            status
        }
    }
`);

const SIGNIN = gql(`
    mutation signin($gameID: ID!, $session: String!, $auth: String!) {
        signin(
            gameID: $gameID,
            session: $session,
            ttl: 9999,
            scope: "0xffffffff",
            authorization: $auth,
        )
    }
`);

const actions = new ethers.utils.Interface([
    "function RESET_MAP() external",
    "function REVEAL_SEED(uint32 blk, uint32 entropy) external",
    "function SPAWN_SEEKER(uint32 sid, uint8 x, uint8 y, uint8 str) external",
    "function MOVE_SEEKER(uint32 sid, uint8 dir) external",
]);

enum BiomeKind {
    UNDISCOVERED,
    BLOCKER,
    GRASS,
    CORN
}

enum Direction {
    NORTH,
    NORTHEAST,
    EAST,
    SOUTHEAST,
    SOUTH,
    SOUTHWEST,
    WEST,
    NORTHWEST
}

const seekers: any = {};

export default class Demo extends Phaser.Scene {

    constructor() {
        super('GameScene');
    }

    preload() {
        this.load.image('tiles', 'assets/roguelikeSheet_transparent.png');
        this.load.spritesheet('chars', 'assets/roguelikeChar_transparent.png', {frameWidth: 16, frameHeight: 16, spacing: 1});
    }

    async create() {
        // which game
        const gameID = "latest";
        const scene = this;

        // setup the client
        const httpLink = new HttpLink({
            uri: COG_HTTP_ENDPOINT,
        });
        const wsLink = new GraphQLWsLink(
            createClient({
                url: COG_WS_ENDPOINT
            }),
        );
        const link = split(
            ({ query }) => {
                const definition = getMainDefinition(query);
                return (
                    definition.kind === 'OperationDefinition' &&
                    definition.operation === 'subscription'
                );
            },
            wsLink,
            httpLink,
        );
        const client = new ApolloClient({
            link,
            uri: COG_HTTP_ENDPOINT,
            cache: new InMemoryCache(),
        });


        // setup wallet providers etc
        const ethereum = (window as any).ethereum;
        if (!ethereum) {
            scene.add.text(100, 100, 'No wallet detected, you need metamask.', { fontSize: '18px' });
            return;
        }
        ethereum.on('accountsChanged', () => {
            console.log('metamask account changed, refreshing...');
            localStorage.clear();
            setTimeout(() => {
                window.location.reload();
            }, 500);
        });
        const provider = new ethers.providers.Web3Provider(ethereum)
        const owner = provider.getSigner();
        try {
            await provider.send("eth_requestAccounts", [])
        } catch {
            scene.add.text(100, 100, 'You must connect your wallet. Refresh to try again', { fontSize: '18px' });
            return;
        }
        const ownerAddr = await owner.getAddress();
        if (localStorage.getItem('ownerAddr') != ownerAddr) {
            localStorage.clear();
            localStorage.setItem('ownerAddr', ownerAddr);
        }

        // setup short lived session key and save in localstorage
        let session: ethers.Wallet;
        let sessionKey = localStorage.getItem('sessionKey');
        if (sessionKey) {
            session = new ethers.Wallet(sessionKey);
            console.log('using session key from localstorage', session.address);
        } else {
            session = ethers.Wallet.createRandom();
            sessionKey = session.privateKey;
            // build signin mutation
            const signin = async () => {
                const msg = ethers.utils.concat([
                    ethers.utils.toUtf8Bytes(`You are signing in with session: `),
                    ethers.utils.getAddress(session.address),
                ]);
                const auth = await owner.signMessage(msg);
                return client.mutate({mutation: SIGNIN, variables: {gameID, auth, session: session.address}});
            }
            // signin with metamask, sign the session key, and save the key for later if success
            const note = scene.add.text(100, 100, 'Check metamask popup for login...', { fontSize: '18px' });
            try {
                await signin();
            } catch {
                scene.add.text(100, 100, 'You must sign in to play. Refresh to try again', { fontSize: '18px' });
                return;
            } finally {
                note.destroy();
            }
            localStorage.setItem('sessionKey', ethers.utils.hexlify(session.privateKey));
        }

        // keep track of seeker owners/sprites
        const getPlayerSeeker = () => {
            for (let k in seekers) {
                if (seekers[k].owner == ownerAddr) {
                    return seekers[k];
                }
            }
            return null;
        }

        // setup dispatch mutation
        const dispatch = async (actionName:string, ...actionArgs:any):Promise<any> => {
            console.log('dispatching', actionName, actionArgs);
            const acts = [actions.encodeFunctionData(actionName, actionArgs)];
            const bundle = ethers.utils.defaultAbiCoder.encode(["bytes[]"], [acts]);
            const actionDigest = ethers.utils.arrayify(ethers.utils.keccak256(bundle));
            console.log('bundle', bundle);
            const auth = await session.signMessage(actionDigest);
            return client.mutate({mutation: DISPATCH, variables: {gameID, auth, actions: acts}})
                .then((res) => {
                    console.log('dispatched', actionName)
                    return res.data.dispatch;
                });
        }

        // plonk dispatch on window so we can call it from the console
        (window as any).dispatch = dispatch;

        // helper to map biomes to tile map index
        const UNDISCOVERED_GRASS = 66;
        const BLOCKING_GRASS = 62;
        const PASSABLE_GRASS = 5;
        const grassType = (i:number, seed:number): number => {
            switch (i) {
                case null: return 1;
                case BiomeKind.UNDISCOVERED: return UNDISCOVERED_GRASS;
                case BiomeKind.BLOCKER: return BLOCKING_GRASS;
                case BiomeKind.GRASS: return PASSABLE_GRASS;
                case BiomeKind.CORN: return PASSABLE_GRASS;
                default: throw new Error(`unknown kind=${i}`);
            }
        }

        // init the map
        const data = Array(32).fill(Array(44).fill(0));
        const map = scene.make.tilemap({ data, tileWidth: 16, tileHeight: 16, width: 64, height: 64 });
        const landTiles = map.addTilesetImage('tiles', undefined, 16, 16, 0, 1);
        const baseLayer = map.createBlankLayer('base', landTiles, 0, 0).setPipeline('Light2D');
        const resourcesLayer = map.createBlankLayer('resources', landTiles, 0, 0).setPipeline('Light2D');
        map.fill(BLOCKING_GRASS, 0, 0, 48,32, undefined, baseLayer);

        // bit of lighting to highlight player
        const light = this.lights.addLight(0, 0, 50).setIntensity(0.15);
        this.lights.enable().setAmbientColor(0xfafafa);

        // UI
        const sideBarDefaults = {
            fixedWidth: 150,
            align: 'justify',
            fontFamily: 'courier',
            color: '#699625',
        };
        const lhs = 536;
        const title = scene.add.text(lhs, 16, 'CORNSEEKERS', { fontSize: '18px', ...sideBarDefaults });
        const playerTitle = scene.add.text(lhs, 125, 'PLAYER', { ...sideBarDefaults, fontSize: '12px' });
        const playerName = scene.add.text(lhs, 140, `Owner: ${ownerAddr.slice(0,16)}`, { ...sideBarDefaults, fontSize: '10px' });
        const playerBalance = scene.add.text(lhs, 155, '', { ...sideBarDefaults, fontSize: '10px' });
        const leaderboardTitle = scene.add.text(lhs, 200, 'LEADERBOARD', { ...sideBarDefaults, fontSize: '12px' });
        const leaders = Array(8).fill(null).map((_, i) => {
            return scene.add.text(lhs, 220+(i*15), `${i+1} ---`, { ...sideBarDefaults, fontSize: '10px' })
        });
        const help = scene.add.text(lhs, 40, [
            `USE W,A,S,D TO MOVE YOUR`,
            `SEEKER. STAND NEAR AN `,
            `UNDISCOVERED AREA TO SCOUT`,
            `WHAT IS INSIDE THE AREA.`,
            `COLLECT CORN BY STANDING `,
            `ON A CORN TILE`,
        ], { fontSize: '9px', ...sideBarDefaults });
        const buttonStyle = {
            ...sideBarDefaults,
            fontSize: '8px',
            color: '#acee44',
            backgroundColor: '#699625',
            padding: {x: 5, y: 5},
        };
        const resetButton = scene.add.text(lhs, 420, 'INIT MAP', buttonStyle)
            .setInteractive()
            .on('pointerup', () => {
                const ans = prompt("Initialize the map.\n\nThis should only be done once, it will act a bit weird if you do it do an already running game.\n\nType 'yes' to do it anyway.");
                if (ans != "yes") {
                    return;
                }
                // reset the map
                dispatch('RESET_MAP')
            });
        const signoutButton = scene.add.text(lhs, 448, 'CLEAR SESSION', buttonStyle)
            .setInteractive()
            .on('pointerup', () => {
                // reset the map
                localStorage.clear();
                (window as any).location.reload();
            });
        const spawnButton = scene.add.text(lhs, 475, 'SPAWN SEEKER', buttonStyle)
            .setInteractive()
            .on('pointerup', () => {
                if (getPlayerSeeker()) {
                    alert('Sorry, you already have a seeker!');
                    return;
                }
                // spawn a seeker at a random location along the top edge
                dispatch('SPAWN_SEEKER', Object.keys(seekers).length+2, Math.floor(Math.random()*32), 0, 1)
            });

        // keep track of things we are trying to reveal
        const revealing = {} as any;

        // helper to check if a coord is "near" the player selection
        const isNearPlayer = (x:number,y:number):boolean => {
            const seeker = getPlayerSeeker();
            if (!seeker) {
                return false;
            }
            return Math.abs(x - (seeker.sprite.x/map.tileWidth)) < 4
                && Math.abs(y - (seeker.sprite.y/map.tileHeight)) < 4;
        }

        // helper to map a seeker id to a char sprite
        const charIdx = (key:string): number => {
            const id = BigNumber.from(key).toNumber();
            return id % 14;
        }

        // this is the main update loop which fires each time a
        // the subscription gets an update of the world state
        const onStateChange = (state: any) => {
            console.log(state);
            if (!state) {
                return;
            }
            const undiscovered = [] as any;

            // draw the map
            state.tiles.forEach((tile: any) => {
                const x = BigNumber.from(tile.coords[0]).toNumber();
                const y = BigNumber.from(tile.coords[1]).toNumber();
                const blk = BigNumber.from(tile.seed?.key || 0).toNumber();
                // we use different types of grass to indicate passable/blocking
                const grass = grassType(tile.biome, blk+x+y);
                map.putTileAt(grass, x, y, undefined, baseLayer);
                // then we place something pretty on top to make the tile distinct
                if (tile.biome === BiomeKind.CORN) {
                    const corn = 15;
                    map.putTileAt(corn, x, y, undefined, resourcesLayer);
                } else if (tile.biome == BiomeKind.BLOCKER) {
                    const tree = [526,592,649,526,526][(blk+x+y) % 5];
                    map.putTileAt(tree, x, y, undefined, resourcesLayer);
                } else if (tile.biome == BiomeKind.UNDISCOVERED) {
                    undiscovered.push({x,y});
                } else {
                    map.removeTileAt(x, y, undefined, undefined, resourcesLayer);
                }
                // resolve any nearby pending tiles that needs resolving
                if (tile.seed && tile.biome == BiomeKind.UNDISCOVERED) {
                    // resolve if nearby
                    if (isNearPlayer(x,y)) {
                        // hint this tile is pending
                        map.putTileAt(BLOCKING_GRASS, x, y, undefined, baseLayer);
                        if (!revealing[`${x}-${y}`]) {
                            // generate some randomness ... obvisouly letting the
                            // client decide random is bad - this is just a toy
                            const entropy = Math.floor(Math.random()*1000);
                            // start reveal
                            revealing[`${x}-${y}`] = true;
                            dispatch("REVEAL_SEED", blk, entropy)
                                .catch((err) => console.error(`REVEAL_SEED ${blk} ${entropy} fail`, err));
                        }
                    }
                }
            });

            // prettify the undiscovered area
            undiscovered.forEach(({x,y}: any) => {
                const tile = map.getTileAt(x,y, undefined, baseLayer);
                const tileAbove = map.getTileAt(x,y-1, undefined, baseLayer)?.index !== UNDISCOVERED_GRASS;
                const tileBelow = map.getTileAt(x,y+1, undefined, baseLayer)?.index !== UNDISCOVERED_GRASS;
                const tileLeft = map.getTileAt(x-1,y, undefined, baseLayer)?.index !== UNDISCOVERED_GRASS;
                const tileRight = map.getTileAt(x+1,y, undefined, baseLayer)?.index !== UNDISCOVERED_GRASS;
                const isRevealing = revealing[`${x}-${y}`];
                const prettyTile = () => {
                    if (tileAbove && tileBelow && tileLeft && tileRight) {
                        return 692;
                    } else if (tileAbove  && tileBelow  && tileLeft ) {
                        return 690;
                    } else if (tileAbove  && tileBelow  && tileRight ) {
                        return 689;
                    } else if (tileBelow  && tileLeft  && tileRight ) {
                        return 632;
                    } else if (tileAbove  && tileLeft  && tileRight ) {
                        return 633;
                    } else if (tileAbove  && tileLeft ) {
                        return 520;
                    } else if (tileAbove  && tileRight ) {
                        return 522;
                    } else if (tileBelow  && tileLeft ) {
                        return 634;
                    } else if (tileBelow  && tileRight ) {
                        return 636;
                    } else if (tileLeft  && tileRight ) {
                       return 408;
                    } else if (tileAbove ) {
                        return 521;
                    } else if (tileBelow ) {
                        return 635;
                    } else if (tileLeft ) {
                        return 577;
                    } else if (tileRight ) {
                        return 579;
                    } else {
                        return 578;
                    }
                }
                map.putTileAt(isRevealing ? 1376 : prettyTile(), x, y, undefined, resourcesLayer);
            })

            // update the seeker locations on the map
            // this is a bit verbose as we animate the movement, ultimately we
            // are simply reading which tile the seeker is positioned at and translating
            // that into the px locations
            state.seekers.forEach((seeker: any) => {
                // find the position of the seekers
                const x = BigNumber.from(seeker.position.coords[0]).toNumber();
                const y = BigNumber.from(seeker.position.coords[1]).toNumber();
                const sx = 16*x+8;
                const sy = 16*y+8;
                if (!seeker.player) {
                    return;
                }
                const owner = ethers.utils.getAddress(seeker.player.address);
                const isPlayerSeeker = owner == ownerAddr;
                if (!seekers[seeker.key]) {
                    const sprite = scene.add.image(16,16,'chars', charIdx(seeker.key));
                    const targets = isPlayerSeeker ? [sprite, light] : [sprite];
                    const tween = scene.tweens.add({
                        targets,
                        x: sx,
                        y: sy,
                        duration: 500,
                        ease: 'Power2',
                        paused: false,
                    });
                    sprite.x = sx;
                    sprite.y = sy;
                    if (isPlayerSeeker) {
                        light.x = sprite.x;
                        light.y = sprite.y;
                    }
                    seekers[seeker.key] = {
                        key: seeker.key,
                        sprite,
                        tween,
                        owner,
                    };
                }
                const {sprite,tween} = seekers[seeker.key];
                if (sprite.x != sx || sprite.y != sy) {
                    if (tween.isPlaying()) {
                        tween.updateTo('x', sx, true);
                        tween.updateTo('y', sy, true);
                    }else{
                        const targets = isPlayerSeeker ? [sprite,light] : [sprite];
                        seekers[seeker.key].tween = scene.tweens.add({
                            targets,
                            x: sx,
                            y: sy,
                            duration: 500,
                            ease: 'Power2',
                            paused: false,
                        });
                    }
                }
                // update player's sidebar stats info
                if (isPlayerSeeker) {
                    // score
                    playerBalance.setText(`Corn collected: ${seeker.cornBalance}`);
                }
            });

            // update the leaderboard
            [...state.seekers].sort((a:any,b:any) => {
                return b.cornBalance - a.cornBalance;
            }).forEach((seeker:any, i:number) => {
                if (i > leaders.length-1) {
                    return;
                }
                if (!seeker.player) {
                    return;
                }
                leaders[i].setText(`${i+1} - ${seeker.player.address.slice(0, 8)} - ${seeker.cornBalance} corns`);
            });


        }

        // helper to check if a tile is passable
        const isBlocker = (tile: any):boolean => {
            if (!tile) {
                return true;
            }
            return tile.index !== PASSABLE_GRASS;
        };

        // track pending moves as light on the tiles
        const pendingMoves:any = [];
        (window as any).pendingMoves = pendingMoves;

        // get where we _think_ the seeker is, either based on
        // the last known pending move, or based on the seeker's pos
        const getExpectedPlayerWorldXY = (): number[] => {
            if (pendingMoves.length > 0) {
                const {tile} = pendingMoves[pendingMoves.length-1];
                return [tile.x*16, tile.y*16];
            }
            const seeker = getPlayerSeeker();
            if (!seeker) {
                console.log('no seeker loc')
                return [-1, -1];
            }
            return [seeker.sprite.x-8, seeker.sprite.y-8];
        }

        //  dispatch MOVE_SEEKER on WASD movement
        const move = (dir: Direction) => async () => {
            if (pendingMoves.length > 10) {
                console.log('soft limit on pending moves');
                return;
            }
            const seeker = getPlayerSeeker();
            if (!seeker) {
                alert("You have no seeker.\n\n Click spawn seeker to start.");
                return;
            }
            // get the position where we think we are based on pending txs
            const [x,y] = getExpectedPlayerWorldXY();
            if (x == -1 || y == -1) {
                console.error('failed to find player position');
                return;
            }
            // check we don't try and move onto a blocker
            let tile;
            switch (dir) {
                case Direction.NORTH:
                    tile = baseLayer.getTileAtWorldXY(x, y+16, true);
                    if (isBlocker(tile)) {
                        console.log('bumped into a wall');
                        return;
                    }
                    break;
                case Direction.SOUTH:
                    tile = baseLayer.getTileAtWorldXY(x, y-16, true);
                    if (isBlocker(tile)) {
                        console.log('bumped into a wall');
                        return;
                    }
                    break;
                case Direction.EAST:
                    tile = baseLayer.getTileAtWorldXY(x+16, y, true);
                    if (isBlocker(tile)) {
                        console.log('bumped into a wall');
                        return;
                    }
                    break;
                case Direction.WEST:
                    tile = baseLayer.getTileAtWorldXY(x-16, y, true);
                    if (isBlocker(tile)) {
                        console.log('bumped into a wall');
                        return;
                    }
                    break;
            }
            if (!tile) {
                return;
            }
            // add a light to indicate pending tx
            const light = scene.lights.addPointLight(tile.x*16+8, tile.y*16+8, 0xffffff, 15, 0.05, 0.07)
            try {
                pendingMoves.push({ light, tile });
                const res = await dispatch("MOVE_SEEKER", seeker.key, dir);
                const mv = pendingMoves.find((p:any) => p.light == light);
                if (mv) {
                    mv.id = res.id;
                } else {
                    console.log('failed to add id to light');
                }
            } catch(err) {
                console.error('moveFail', err);
                light.destroy();
                const mv = pendingMoves.find((p:any) => p.light == light);
                if (mv) {
                    pendingMoves.splice(pendingMoves.indexOf(mv), 1);
                }
            }
        }
        scene.input.keyboard.on('keydown-A', move(Direction.WEST));
        scene.input.keyboard.on('keydown-D', move(Direction.EAST));
        scene.input.keyboard.on('keydown-W', move(Direction.SOUTH));
        scene.input.keyboard.on('keydown-S', move(Direction.NORTH));

        const onTxChange = async (tx:any) => {
            for (let i=0; i<pendingMoves.length; i++) {
                const {light, tile, id} = pendingMoves[i];
                if (!id) {
                    continue;
                }
                if (tx.id == id && tx.status != "PENDING") {
                    for (let j=0; j<i+1; j++) {
                        pendingMoves[j].light.destroy();
                    }
                    pendingMoves.splice(0, i+1);
                    return;
                }
            }
        }

        // subscribe to future state changes
        client.subscribe({
            query: STATE_SUBSCRIPTION,
        }).subscribe(
            (result) => onStateChange(result.data.state),
            (err) => console.error('stateSubscriptionError', err),
            () => console.warn('stateSubscriptionClosed')
        );


        // watch for tx belonging to this player
        client.subscribe({
            query: TX_SUBSCRIPTION,
            variables: {owner: ownerAddr},
        }).subscribe(
            (result) => onTxChange(result.data.transaction),
            (err) => console.error('txSubscriptionError', err),
            () => console.warn('txSubscriptionClosed')
        )

        // fetch initial state
        await client.query({query: STATE_QUERY})
            .then((result) => onStateChange(result.data.game.state))
            .catch((err) => console.error('err', err));

    }

}
