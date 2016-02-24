import { bindActionCreators, Dispatch } from 'redux';
import { connect } from 'react-redux';
import * as React from 'react';

import Header from '../components/Header';
import MainSection from '../components/MainSection';
import * as TodoActions from '../actions/todos';
import { File } from '../models/files';
import { Todo } from '../models/todos';

const style = require('../styles/app.styl')

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
      <div className={style.app}>
        <Header
          file={file}
          addTodo={actions.addTodo} />
        <MainSection
          todos={todos}
          actions={actions} />
      </div>
    );
  }
}

const mapStateToProps = state => ({
  todos: state.todos,
  file: state.file
});

export default connect(mapStateToProps)(App);
