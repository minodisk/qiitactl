var LiveReloadPlugin = require('webpack-livereload-plugin')

module.exports = {
  entry: './src/main.tsx',
  output: {
    filename: 'bundle.js'
  },
  resolve: {
    extensions: [
      '',
      '.ts', '.tsx', '.js',
      '.styl', '.css'
    ]
  },
  module: {
    loaders: [
      { test: /\.tsx?$/, loaders: ['ts'] },
      {
        test: /\.css/,
        loaders: [
          'style?sourceMap',
          'css'
        ]
      },
      {
        test: /\.s[ac]ss/,
        loaders: [
          'style',
          'css?sourceMap',
          'sass?sourceMap'
        ]
      },
      { test: /\.styl/, loaders: ['style', 'css', 'stylus'] },
      { test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/, loader: "url-loader?limit=10000&mimetype=application/font-woff" },
      { test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/, loader: "file-loader" }
    ]
  },
  plugins: [
    new LiveReloadPlugin()
  ]
}
