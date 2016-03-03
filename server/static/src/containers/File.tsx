import * as React from 'react';
import { connect } from 'react-redux';

import { File } from '../models/file'
import {
  getFile,
  watchFile,
  unwatchFile
} from '../actions/file'

const styles = require('../styles/content.css')

interface Props {
  location: any;
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
    if (nextProps.location.pathname == this.props.location.pathname) {
      return
    }
    nextProps.watchFile()
  }

  componentWillUnmount() {
    this.props.unwatchFile()
  }

  render() {
    console.log(this.props.file)
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
