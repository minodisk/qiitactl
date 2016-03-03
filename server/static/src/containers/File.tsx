import * as React from 'react';
import { connect } from 'react-redux';

import {
  getFile,
  watchFile,
  unwatchFile
} from '../actions/file'

const styles = require('../styles/content.css')

interface Props {
  location: any;
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
    nextProps.watchFile()
  }

  componentWillUnmount() {
    this.props.unwatchFile()
  }

  render() {
    return (
      <p>{this.props.location.pathname}</p>
    );
  }
}

const ejectRootPath = path => path.replace(/^\/markdown\//, '')

const mapStateToProps = state => ({
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
