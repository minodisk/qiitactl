import * as React from 'react';

const styles = require('../styles/content.css')

interface Props {};
interface State {};

export default class NotFound extends React.Component<Props, State> {
  render() {
    console.log('NotFound.render')
    return (
      <p>NotFound</p>
    );
  }
}
