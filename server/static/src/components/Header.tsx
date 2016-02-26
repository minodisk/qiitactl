import * as React from 'react';

import TodoTextInput from './TodoTextInput';
import File from './File'
import * as model from '../models/files'

const style = require('../styles/header.css')

interface HeaderProps {
  file: model.File;
  addTodo: Function;
};

export default class Header extends React.Component<HeaderProps, void> {
  handleSave(text) {
    if (text.length !== 0) {
      this.props.addTodo(text);
    }
  }

  render() {
    const {file} = this.props
    return (
      <header className={style.header}>
        <nav>
          {(() => {
            if (file != null) {
              return (
                <File
                  indent={0}
                  file={file}
                  />
              )
            }
          })()}
        </nav>
        <TodoTextInput
          newTodo
          onSave={this.handleSave.bind(this)}
          placeholder="What needs to be done?" />
      </header>
    );
  }
}
