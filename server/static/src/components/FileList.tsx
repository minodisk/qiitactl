import * as React from 'react';
const LeftNav = require('material-ui/lib/left-nav')
const Paper = require('material-ui/lib/paper')

import TodoTextInput from './TodoTextInput';
import File from './File'
import * as model from '../models/files'

const {leftNav, leftNavInner} = require('../styles/file-list.css')

interface Props {
  file: model.File;
  addTodo: Function;
};

interface State {
  open: boolean
}

export default class FileList extends React.Component<Props, State> {
  constructor(props) {
    super(props)
    this.state = {open: true}
  }

  render() {
    const {file} = this.props
    return (
      <LeftNav
        className={leftNav}
        open={this.state.open}
      >
        <Paper className={leftNavInner}>
          <File
            indent={0}
            file={file}
            />
        </Paper>
      </LeftNav>
    );
  }
}
