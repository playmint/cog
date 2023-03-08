module.exports = {
    root: true,
    overrides: [
        {
            files: ["*.js"],
            processor: "@graphql-eslint/graphql",
            extends: ["eslint:recommended", "plugin:prettier/recommended"],
            env: {
                node: true,
                es6: true,
            },
        },
        {
            files: ["*.graphql", "*.graphqls"],
            parser: "@graphql-eslint/eslint-plugin",
            plugins: ["@graphql-eslint"],
            rules: {
                "prettier/prettier": "error",
            },
        },
        // the following is required for `eslint-plugin-prettier@<=3.4.0` temporarily
        // after https://github.com/prettier/eslint-plugin-prettier/pull/415
        // been merged and released, it can be deleted safely
        {
            files: ["*.js/*.graphql", "*.js/*.graphqls"],
            rules: {
                "prettier/prettier": "off",
            },
        },
    ],
};
