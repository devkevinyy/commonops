
const path = require('path');
const webpack = require("webpack");

module.exports = {
    // mode: 'development',
    devtool: 'false',
    mode: "production",
    entry: {
        main: path.resolve(__dirname, './src/index.js'),
        vendor: ['react','react-dom','react-router-dom', 'axios'],
    },
    output: {
        path: path.resolve(__dirname, './public/'),
        filename: '[name].dll.js',
        library: '[name]_library'
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/i,
                exclude: /node_modules/,
                loader: "babel-loader",
            },
            {
                test: /\.(css|less)$/,
                loader: 'style-loader!css-loader'
            },
            {
                test: /\.png$/,
                loader: 'url-loader'
            }
        ]
    },
    plugins: [
        new webpack.optimize.ModuleConcatenationPlugin(),
    ],
    optimization: {
        splitChunks: {
            name: 'vendor',
        }
    }
};