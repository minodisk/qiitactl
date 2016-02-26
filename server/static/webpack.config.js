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
      require('autoprefixer'),
      require('precss'),
      require('postcss-url')
    ]
  },
  plugins: [
    new LiveReloadPlugin()
  ]
}
