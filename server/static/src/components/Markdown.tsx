import * as React from 'react';

import {
  watchFile,
  unwatchFile
} from '../actions/files'

const styles = require('../styles/content.css')

interface Props {
  location: any;
};
interface State {};

export default class Markdown extends React.Component<Props, State> {
  componentWillMount() {
    watchFile(this.props.location.pathname)
  }

  componentWillReceiveProps(nextProps) {
    unwatchFile(this.props.location.pathname)
    watchFile(nextProps.location.pathname)
  }

  componentWillUnmount() {
    console.log('componentWillUnmount')
    unwatchFile(this.props.location.pathname)
  }

  render() {
    return (
      <p>{this.props.location.pathname}</p>
    );
  }
}
