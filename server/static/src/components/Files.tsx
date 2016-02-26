import * as React from 'react';
import {findDOMNode} from 'react-dom'

import File from './File'
import * as models from '../models/files'
import * as classnames from 'classnames'

const {filesOpened, filesClosed} = require('../styles/file.css')

interface FilesProps {
  indent: number;
  files: models.File[]
  opened: boolean,
}

export default class Files extends React.Component<FilesProps, void> {
  render() {
    const {indent, files, opened} = this.props
    return (
      <ul
        className={opened ? filesOpened : filesClosed}
      >
        {
          files.map((file) => <li key={file.id}><File indent={indent} file={file} /></li>)
        }
      </ul>
    )
  }
}
