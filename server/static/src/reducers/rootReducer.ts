import { combineReducers } from 'redux';
import { routerStateReducer } from 'redux-router'

import todos from './todos';
import tree from './tree';
import file from './file';

const rootReducer = combineReducers({
  router: routerStateReducer,
  todos: todos,
  tree: tree,
  file: file,
});

export { rootReducer };
