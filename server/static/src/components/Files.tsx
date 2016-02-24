import * as React from 'react';

import File from './File'
import * as models from '../models/files'

const style = require('../styles/header.styl')

interface FilesProps {
  files: models.File[]
  opened: boolean
}

export default class Files extends React.Component<FilesProps, void> {
  render() {
    const {files, opened} = this.props
    if (!opened) {
      return <div>empty</div>
    }
    return (
      <ul className={style.files}>
      {
        files.map((file) => <File file={file} />)
      }
      </ul>
    )
  }
}
