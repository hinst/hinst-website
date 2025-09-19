//webpack.config.js
const path = require('path');

module.exports = {
	mode: "development",
	devtool: "inline-source-map",
	entry: {
		initialize: "./pages/TypeScript/initialize.ts",
		main: "./pages/TypeScript/main.ts",
	},
	output: {
		path: path.resolve(__dirname, './pages/static/js/dist'),
		filename: "[name].js" // <--- Will be compiled to this single file
	},
	resolve: {
		extensions: [".ts", ".tsx", ".js"],
	},
	module: {
		rules: [
			{
				test: /\.tsx?$/,
				loader: "ts-loader"
			}
		]
	}
};
