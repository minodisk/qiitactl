/// <reference path="../../typings/tsd.d.ts" />

import { assign } from 'lodash';
import { handleActions, Action } from 'redux-actions';

import { Todo } from '../models/todos';
import {
  UPDATE_FILE_LIST,
} from '../actions/files';

const initialState = [<Todo>{
  text: 'Use Redux with TypeScript',
  completed: false,
  id: 0
}];

export default handleActions<Todo[]>({
  [ADD_TODO]: (state: Todo[], action: Action): Todo[] => {
    return [{
      id: state.reduce((maxId, todo) => Math.max(todo.id, maxId), -1) + 1,
      completed: action.payload.completed,
      text: action.payload.text
    }, ...state];
  },

  [COMPLETE_ALL]: (state: Todo[], action: Action): Todo[] => {
    const areAllMarked = state.every(todo => todo.completed);
    return <Todo[]>state.map(todo => assign({}, todo, {
      completed: !areAllMarked
    }));
  },

  [CLEAR_COMPLETED]: (state: Todo[], action: Action): Todo[] => {
    return state.filter(todo => todo.completed === false);
  }
}, initialState);
