import * as React from 'react';

import Files from './Files'
import Link from './Link'
import * as models from '../models/files'

interface FileProps {
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
    const {file} = this.props
    if (file.children == null) {
      return (
        <li>
          <Link
            active={true}
            children={file.name}
            onClick={() => console.log('start')} />
        </li>
      )
    }
    return (
      <li>
        <Link
          active={true}
          children={file.name}
          onClick={() => this.toggleOpen()} />
        <Files
          files={file.children}
          opened={this.state.opened} />
      </li>
    )
  }
}
