import { handleActions, Action } from 'redux-actions';

import { File } from '../models/file';
import {
  DID_GET_FILE,
  DID_CHANGE_FILE,
} from '../constants/ActionTypes';

export default handleActions<File>({
  [DID_GET_FILE]: (state: File, action: Action): File => {
    return action.payload;
  },
  [DID_CHANGE_FILE]: (state: File, action: Action): File => {
    return action.payload;
  },
}, null);
