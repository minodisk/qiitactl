import { combineReducers } from 'redux';

import todos from './todos';
import file from './file';

const rootReducer = combineReducers({
  todos: todos,
  file: file
});

export { rootReducer };
