import * as React from 'react';

import Files from './Files'
import Element from './Element'
import * as models from '../models/files'

const styles = require('../styles/file.css')

interface Props {
  indent: number;
  file: models.File
}

interface State {
  opened: boolean;
}

export default class File extends React.Component<Props, State> {
  constructor(props, context) {
    super(props, context)
    this.state = {
      opened: true,
    }
  }

  toggleOpen() {
    const {opened} = this.state
    this.setState({
      opened: !opened,
    })
  }

  render() {
    const {indent, file} = this.props
    if (file.children == null) {
      return (
        <div>
          <Element
            indent={indent}
            className={styles.file}
            active={true}
            title={file.abs}
            children={file.name}
            linkTo={file.rel}
          />
        </div>
      )
    }
    return (
      <div>
        <Element
          indent={indent}
          className={this.state.opened ? styles.dirOpened : styles.dirClosed}
          active={true}
          title={file.abs}
          children={file.name}
          onClick={() => this.toggleOpen()}
        />
        <Files
          indent={indent + 1}
          files={file.children}
          opened={this.state.opened}
        />
      </div>
    )
  }

}
