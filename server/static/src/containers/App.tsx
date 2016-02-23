/// <reference path="../../typings/tsd.d.ts" />

import {bindActionCreators, Dispatch} from 'redux'
import {connect} from 'react-redux'
import * as React from 'react'
import FileList from '../components/FileList'
import * as FilesActions from '../actions/files'

interface AppProps {
  paths: string[];
  dispatch: Dispatch;
}

class App extends React.Component<AppProps, void> {
  render() {
    const {paths, dispatch} = this.props
    const actions = bindActionCreators(FilesActions, dispatch)
    return (
      <div>
        <FileList
          paths={paths}
          actions={actions}
        />
      </div>
    )
  }
}

const mapStateToProps = state => ({
  paths: state.paths
})

export default connect(mapStateToProps)(App)
