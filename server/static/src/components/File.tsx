import * as React from 'react';

import Files from './Files'
import Link from './Link'
import * as models from '../models/files'

const classes = require('../styles/header.styl')

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
        <div className={classes.file}>
          <Link
            active={true}
            children={file.name}
            onClick={() => console.log('start')} />
        </div>
      )
    }
    return (
      <div className={classes.dir}>
        <Link
          active={true}
          children={file.name}
          onClick={() => this.toggleOpen()} />
        <Files
          files={file.children}
          opened={this.state.opened} />
      </div>
    )
  }
}
