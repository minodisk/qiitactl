import { assign } from 'lodash';
import { handleActions, Action } from 'redux-actions';

import { File } from '../models/files';
import {
  REQUEST_FILES,
  RECIEVE_FILES,
  // WILL_WATCH_FILE,
  // DID_WATCH_FILE,
  // CHANGE_FILE,
} from '../constants/ActionTypes';

export default handleActions<File>({
  [REQUEST_FILES]: (state: File, action: Action): File => {
    return state
  },
  [RECIEVE_FILES]: (state: File, action: Action): File => {
    return action.payload;
  },
  // [WILL_WATCH_FILE]: (state: File, action: Action): File => {
  //   console.log('WILL_WATCH_FILE')
  //   return state
  // },
  // [DID_WATCH_FILE]: (state: File, action: Action): File => {
  //   console.log('DID_WATCH_FILE')
  //   return state
  // },
  // [CHANGE_FILE]: (state: File, action: Action): File => {
  //   console.log('CHANGE_FILE')
  //   return state
  // },
}, <File>{name: 'loading...', children: []});
