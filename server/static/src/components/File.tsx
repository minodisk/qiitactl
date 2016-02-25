import * as React from 'react';

import Files from './Files'
import Link from './Link'
import * as models from '../models/files'

const styles = require('../styles/file.scss')

interface FileProps {
  indent: number;
  file: models.File
}

interface FileState {
  opened: boolean;
}

export default class File extends React.Component<FileProps, FileState> {
  constructor(props, context) {
    super(props, context)
    this.state = {
      opened: true
    }
  }

  toggleOpen() {
    const {opened} = this.state
    this.setState({opened: !opened})
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
            children={file.name}
            onClick={() => console.log('start')} />
        </div>
      )
    }
    return (
      <div>
        <Link
          indent={indent}
          className={this.state.opened ? styles.dirOpened : styles.dirClosed}
          active={true}
          children={file.name}
          onClick={() => this.toggleOpen()} />
        <Files
          indent={indent + 1}
          files={file.children}
          opened={this.state.opened} />
      </div>
    )
  }
}
