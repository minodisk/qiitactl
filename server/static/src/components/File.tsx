import * as React from 'react';

import Files from './Files'
import Link from './Link'
import * as models from '../models/files'

const styles = require('../styles/file.css')

interface Props {
  indent: number;
  file: models.File
}

interface State {
  opened: boolean;
  showing: boolean;
}

export default class File extends React.Component<Props, State> {
  constructor(props, context) {
    super(props, context)
    this.state = {
      opened: true,
      showing: false
    }
  }

  toggleOpen() {
    const {opened} = this.state
    this.setState({
      opened: !opened,
      showing: this.state.showing
    })
  }

  render() {
    const {indent, file} = this.props
    if (file.children == null) {
      return (
        <div>
          <Link
            indent={indent}
            className={styles.file}
            active={true}
            title={file.path}
            children={file.name}
            onClick={this.handleClick} />
        </div>
      )
    }
    return (
      <div>
        <Link
          indent={indent}
          className={this.state.opened ? styles.dirOpened : styles.dirClosed}
          active={true}
          title={file.path}
          children={file.name}
          onClick={() => this.toggleOpen()} />
        <Files
          indent={indent + 1}
          files={file.children}
          opened={this.state.opened} />
      </div>
    )
  }

  handleClick = (e) => {
    this.setState({
      opened: this.state.opened,
      showing: true
    })
  }
}
