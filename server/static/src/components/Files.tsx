import * as React from 'react';

import File from './File'
import * as models from '../models/files'
import * as classnames from 'classnames'

const style = require('../styles/header.styl')

interface FilesProps {
  files: models.File[]
  opened: boolean
}

export default class Files extends React.Component<FilesProps, void> {
  render() {
    const {files, opened} = this.props
    return (
      <ul className={classnames({
        [style.files]: true,
        [style.closed]: !opened
      })}>
      {
        files.map((file) => <li><File file={file} /></li>)
      }
      </ul>
    )
  }
}
