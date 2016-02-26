import * as React from 'react';

import { Todo } from '../models/todos';
import TodoItem from './TodoItem';
import {
  SHOW_ALL,
  SHOW_COMPLETED,
  SHOW_ACTIVE
} from '../constants/TodoFilters';

const styles = require('../styles/content.css')

const TODO_FILTERS = {
  [SHOW_ALL]: () => true,
  [SHOW_ACTIVE]: todo => !todo.completed,
  [SHOW_COMPLETED]: todo => todo.completed
};

interface Props {
  todos: Todo[];
  actions: any;
};
interface State {
  filter: string;
};

export default class Content extends React.Component<Props, State> {
  constructor(props, context) {
    super(props, context);
    this.state = { filter: SHOW_ALL };
  }

  handleClearCompleted() {
    const atLeastOneCompleted = this.props.todos.some(todo => todo.completed);
    if (atLeastOneCompleted) {
      this.props.actions.clearCompleted();
    }
  }

  handleShow(filter) {
    this.setState({ filter });
  }

  renderToggleAll(completedCount) {
    const { todos, actions } = this.props;
    if (todos.length > 0) {
      return (
        <input className="toggle-all"
               type="checkbox"
               checked={completedCount === todos.length}
               onChange={() => actions.completeAll()} />
      );
    }
  }

  render() {
    const { todos, actions } = this.props;
    const { filter } = this.state;

    const filteredTodos = todos.filter(TODO_FILTERS[filter]);
    const completedCount = todos.reduce((count: number, todo): number =>
      todo.completed ? count + 1 : count,
      0
    );

    return (
      <section className={styles.content}>
        <div className={styles.contentInner}>
          {this.renderToggleAll(completedCount)}
          <ul className="todo-list">
            {filteredTodos.map(todo =>
              <TodoItem
                key={todo.id}
                todo={todo}
                { ...actions }/>
            )}
          </ul>
        </div>
      </section>
    );
  }
}
