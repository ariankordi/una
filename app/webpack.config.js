const VueLoaderPlugin = require("vue-loader/lib/plugin");

const UglifyJsPlugin = require("uglifyjs-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const OptimizeCSSAssetsPlugin = require("optimize-css-assets-webpack-plugin");

module.exports = {
  entry: ['./main.js', './main.scss'],
	output: {
		filename: '../../static/app.js'
	},
  mode: 'development',
	watch: true,
	module: {
		rules: [{
      test: /\.vue$/,
      loader: 'vue-loader'
    },
    {
      test: /\.js$/,
      loader: 'babel-loader',
      exclude: /node_modules/,
      query: {
        presets: ['es2015']
      }
    },
    {
			test: /\.scss$/,
			use: [
				MiniCssExtractPlugin.loader, "css-loader", "sass-loader",
			]
		}]
	},
  optimization: {
    minimizer: [
      new UglifyJsPlugin({
        cache: true,
        parallel: true,
        sourceMap: true
      }),
      new OptimizeCSSAssetsPlugin({})
    ]
  },
	plugins: [
    new VueLoaderPlugin(),
		new MiniCssExtractPlugin({
			filename: "../../static/app.css",
      chunkFilename: "../../static/app-[id].css"
		})
	]
}
