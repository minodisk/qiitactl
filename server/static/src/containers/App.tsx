import { bindActionCreators, Dispatch } from 'redux';
import { connect } from 'react-redux';
import * as React from 'react';
const AppBar = require('material-ui/lib/app-bar')
const IconButton = require('material-ui/lib/icon-button')
const IconMenu = require('material-ui/lib/menus/icon-menu')
const MenuItem = require('material-ui/lib/menus/menu-item')
const NavigationClose = require('material-ui/lib/svg-icons/navigation/close')
const MoreVertIcon = require('material-ui/lib/svg-icons/navigation/more-vert')
const injectTouchTapEvent = require('react-tap-event-plugin');
injectTouchTapEvent();

import FileList from '../components/FileList';
import Content from '../components/Content';
import * as TodoActions from '../actions/todos';
import { File } from '../models/files';
import { Todo } from '../models/todos';

const styles = require('../styles/app.css')

interface AppProps {
  todos?: Todo[];
  file?: File;
  dispatch?: Dispatch;
}

class App extends React.Component<AppProps, void> {
  render() {
    const { todos, file, dispatch } = this.props;
    const actions = bindActionCreators(TodoActions, dispatch);
    return (
      <div>
        <AppBar
          title='qiitactl'
        />
        <FileList
          file={file}
          addTodo={actions.addTodo}
        />
        <Content
          todos={todos}
          actions={actions}
        />
      </div>
    );
  }
}

const mapStateToProps = state => ({
  todos: state.todos,
  file: state.file
});

export default connect(mapStateToProps)(App);
