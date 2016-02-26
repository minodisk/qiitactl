var LiveReloadPlugin = require('webpack-livereload-plugin')

module.exports = {
  entry: './src/main.tsx',
  output: {
    filename: 'bundle.js',
    publicPath: 'http://localhost:9000/'
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
        loader: 'url-loader?limit=10000&mimetype=application/font-woff'
      },
      {
        test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
        loader: 'file-loader'
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
