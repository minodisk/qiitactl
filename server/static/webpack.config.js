var LiveReloadPlugin = require('webpack-livereload-plugin')

module.exports = {
  entry: './src/main.tsx',
  output: {
    filename: 'bundle.js'
  },
  resolve: {
    extensions: [
      '',
      '.ts', '.tsx', '.js'
    ]
  },
  module: {
    loaders: [
      { test: /\.styl/, loaders: ['style', 'css', 'stylus'] },
      { test: /\.tsx?$/, loaders: ['ts'] }
    ]
  },
  plugins: [
    new LiveReloadPlugin()
  ]
}
