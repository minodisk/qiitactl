import { assign } from 'lodash';
import { handleActions, Action } from 'redux-actions';

import { Tree } from '../models/tree';
import {
  WILL_GET_TREE,
  DID_GET_TREE,
} from '../constants/ActionTypes';

export default handleActions<Tree>({
  [WILL_GET_TREE]: (state: Tree, action: Action): Tree => {
    return state
  },
  [DID_GET_TREE]: (state: Tree, action: Action): Tree => {
    return action.payload;
  },
  // [WILL_WATCH_FILE]: (state: Tree, action: Action): Tree => {
  //   console.log('WILL_WATCH_FILE')
  //   return state
  // },
  // [DID_WATCH_FILE]: (state: Tree, action: Action): Tree => {
  //   console.log('DID_WATCH_FILE')
  //   return state
  // },
  // [CHANGE_FILE]: (state: Tree, action: Action): Tree => {
  //   console.log('CHANGE_FILE')
  //   return state
  // },
}, <Tree>{name: 'loading...', children: []});
