/// <reference path="../../typings/tsd.d.ts" />

import * as React from 'react'
import File from './File'

interface FileListProps {
  paths: string[];
  actions: any;
}
interface FileListState {}

class FileList extends React.Component<FileListProps, FileListState> {
  render() {
    const {paths, actions} = this.props
    return (
      <ul>
        {paths.map((path) =>
          <File
            path={path}
            {...actions}
          />
        )}
      </ul>
    )
  }
}

export default FileList
