import { assign } from 'lodash';
import { handleActions, Action } from 'redux-actions';

import { File } from '../models/files';
import {
  REQUEST_FILES,
  RECIEVE_FILES
} from '../constants/ActionTypes';

export default handleActions<File>({
  [REQUEST_FILES]: (state: File, action: Action): File => {
    return state
  },
  [RECIEVE_FILES]: (state: File, action: Action): File => {
    return action.payload;
  }
}, <File>{});
