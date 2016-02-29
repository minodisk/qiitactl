var path = require('path')
var LiveReloadPlugin = require('webpack-livereload-plugin')
var tsconfig = require('./tsconfig.json')
var entry = tsconfig.files.filter(function (file) {
  return path.basename(file, path.extname(file)) === 'index'
})[0]

module.exports = {
  entry: entry,
  output: {
    path: 'dist/assets',
    filename: 'bundle.js',
    publicPath: 'http://localhost:9000/assets/'
  },
  devtool: 'source-map',
  resolve: {
    extensions: [
      '',
      '.ts', '.tsx', '.js',
      '.styl', '.css'
    ]
  },
  module: {
    loaders: [
      {
        test: /\.tsx?$/,
        loaders: [
          'ts'
        ]
      },
      {
        test: /\.css/,
        loaders: [
          'style',
          'css?sourceMap',
          'postcss',
        ]
      },
      {
        test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
        loader: 'url-loader?name=[sha512:hash:base64:7].[ext]&limit=10000&mimetype=application/font-woff'
      },
      {
        test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
        loader: 'file-loader?name=[sha512:hash:base64:7].[ext]'
      }
    ]
  },
  postcss: function () {
    return [
      require('postcss-modules-local-by-default'),
      require('autoprefixer'),

      // require('postcss-partial-import'),
      // require('postcss-mixins'),
      // require('postcss-advanced-variables'),
      // require('postcss-custom-media'),
      // require('postcss-custom-properties'),
      // require('postcss-media-minmax'),
      // require('postcss-color-function'),
      // require('postcss-nesting'),
      // require('postcss-nested'),
      // require('postcss-custom-selectors'),
      // require('postcss-atroot'),
      // require('postcss-property-lookup'),
      // require('postcss-extend'),
      // require('postcss-selector-matches'),
      // require('postcss-selector-not'),
      require('precss'),

      require('postcss-url'),
    ]
  },
  plugins: [
    new LiveReloadPlugin()
  ]
}
