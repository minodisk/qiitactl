import * as React from 'react';
import { connect } from 'react-redux';

import {
  watchFile,
  unwatchFile
} from '../actions/files'

const styles = require('../styles/content.css')

interface Props {
  location: any;
  startWatching: Function;
  stopWatching: Function;
};
interface State {};

class Markdown extends React.Component<Props, State> {
  componentWillMount() {
    this.props.startWatching()
  }

  componentWillReceiveProps(nextProps) {
    nextProps.startWatching()
  }

  componentWillUnmount() {
    this.props.stopWatching()
  }

  render() {
    return (
      <p>{this.props.location.pathname}</p>
    );
  }
}

const mapStateToProps = state => ({
})

const mapDispatchToProps = (dispatch, props) => ({
  startWatching: () => {
    dispatch(watchFile(props.location.pathname))
  },
  stopWatching: () => {
    dispatch(unwatchFile(props.location.pathname))
  },
})

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Markdown)
