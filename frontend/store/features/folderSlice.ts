import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface FolderState {
  // Map of folderId -> array of conversation IDs
  assignments: Record<string, string[]>;
}

const initialState: FolderState = {
  assignments: {
    'work': [],
    'friends': [],
    'archive': []
  }
};

const folderSlice = createSlice({
  name: 'folders',
  initialState,
  reducers: {
    assignToFolder: (state, action: PayloadAction<{ folderId: string; conversationId: string }>) => {
      const { folderId, conversationId } = action.payload;

      // Initialize array if it doesn't exist (though it should from initial state)
      if (!state.assignments[folderId]) {
        state.assignments[folderId] = [];
      }

      // Avoid duplicates
      if (!state.assignments[folderId].includes(conversationId)) {
        state.assignments[folderId].push(conversationId);
      }
    },
    removeFromFolder: (state, action: PayloadAction<{ folderId: string; conversationId: string }>) => {
      const { folderId, conversationId } = action.payload;
      if (state.assignments[folderId]) {
        state.assignments[folderId] = state.assignments[folderId].filter(id => id !== conversationId);
      }
    },
    toggleFolderAssignment: (state, action: PayloadAction<{ folderId: string; conversationId: string }>) => {
      const { folderId, conversationId } = action.payload;
      if (!state.assignments[folderId]) {
         state.assignments[folderId] = [];
      }

      const index = state.assignments[folderId].indexOf(conversationId);
      if (index === -1) {
        state.assignments[folderId].push(conversationId);
      } else {
        state.assignments[folderId].splice(index, 1);
      }
    }
  },
});

export const { assignToFolder, removeFromFolder, toggleFolderAssignment } = folderSlice.actions;
export default folderSlice.reducer;
