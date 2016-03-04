import * as React from 'react';
import { connect } from 'react-redux';

import { Connection } from '../models/socket'
import { File } from '../models/file'
import {
  getFile,
  watchFile,
  unwatchFile
} from '../actions/file'

const styles = require('../styles/content.css')

interface Props {
  location: any;
  connection: Connection;
  file: File;
  getFile: Function;
  watchFile: Function;
  unwatchFile: Function;
};
interface State {};

class Markdown extends React.Component<Props, State> {
  componentWillMount() {
    this.props.watchFile()
  }

  componentWillReceiveProps(nextProps) {
    if (nextProps.location.pathname != this.props.location.pathname) {
      nextProps.watchFile()
    }
    console.log(this.props.connection.opened, '->', nextProps.connection.opened);
    if (nextProps.connection.opened != this.props.connection.opened) {
      if (nextProps.connection.opened) {
        nextProps.watchFile()
      }
    }
  }

  componentWillUnmount() {
    this.props.unwatchFile()
  }

  render() {
    if (this.props.file == null) {
      return (<pre></pre>)
    }
    return (
      <pre>{this.props.file.content}</pre>
    );
  }
}

const ejectRootPath = path => path.replace(/^\/markdown\//, '')

const mapStateToProps = state => ({
  file: state.file,
  connection: state.connection,
})

const mapDispatchToProps = (dispatch, props) => ({
  watchFile: () => {
    const pathname = ejectRootPath(props.location.pathname)
    dispatch(getFile(pathname))
    dispatch(watchFile(pathname))
  },
  unwatchFile: () => {
    const pathname = ejectRootPath(props.location.pathname)
    dispatch(unwatchFile(pathname))
  },
})

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Markdown)
