/// <reference path="../../typings/tsd.d.ts" />

import * as React from 'react'
import Link from './Link'

interface FileProps {
  path: string;
  actions: any;
}

interface FileState {}

class File extends React.Component<FileProps, FileState> {
  render() {
    const {path, actions} = this.props
    return (
      <li>
        <Link
          path={path}
          {...actions}
        />
      </li>
    );
  }
}

export default File
