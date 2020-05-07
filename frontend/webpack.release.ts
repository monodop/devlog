import * as webpack from "webpack";
import * as HtmlWebPackPlugin from "html-webpack-plugin";
import * as path from 'path';

const htmlPlugin = new HtmlWebPackPlugin({
    template: "./src/index.html"
});

const config: webpack.Configuration = {
    mode: "production",
    entry: "./src/index.tsx",
    resolve: {
        extensions: [".ts", ".tsx", ".js", ".json"]
    },

    module: {
        rules: [
            { test: /\.(png|svg|jpg|gif)$/, loader: 'file-loader' },
            { test: /\.tsx?$/, exclude: /node_modules/, loader: "awesome-typescript-loader" },
        ]
    },
    
    plugins: [htmlPlugin],
};

export default config;