import React, { FC } from 'react';
import intl from 'react-intl-universal';
import { queryCache, useMutation } from 'react-query';
import { css } from 'emotion';
import { Modal, Button, ConfirmModal } from '@grafana/ui';

import { deleteReportGroup, updateReportGroup } from 'api';

import { useOptimisticMutation, useToggle } from 'hooks';
import { ReportGroupKey } from '../../common/enums';
import { ReportGroupMemberList } from './ReportGroupMemberList';
import { EditReportGroupForm } from './EditReportGroupForm';
import { ReportGroup } from 'common/types';

type Props = {
  onClose: () => void;
  isOpen: boolean;
  reportGroup: ReportGroup;
};

const modalAdjustments = css`
  top: 0;
  bottom: 0;
  width: 80%;
`;

export const EditReportGroupModal: FC<Props> = ({ reportGroup, onClose, isOpen }) => {
  const [deleteAlertIsOpen, setDeleteAlertIsOpen] = useToggle(false);

  const [updateGroup] = useOptimisticMutation<ReportGroup[], ReportGroup, ReportGroup, ReportGroup[]>(
    ['reportGroup'],
    updateReportGroup,
    group => group,
    (prevState, optimisticValue) => {
      if (prevState) {
        const idx = prevState.findIndex(group => group.id === optimisticValue.id);
        if (idx >= 0) {
          prevState[idx] = optimisticValue;
          return [...prevState];
        }
      }
      return prevState;
    },
    []
  );

  const [deleteGroup] = useMutation(deleteReportGroup, {
    onSuccess: () => queryCache.refetchQueries(['reportGroup']),
  });

  const onConfirmDeleteGroup = () => {
    deleteGroup(reportGroup);
    setDeleteAlertIsOpen();
    onClose();
  };

  const onUpdateReportGroup = (key: ReportGroupKey, newValue: string) => {
    const newState: ReportGroup = { ...reportGroup, [key]: newValue };
    updateGroup(newState);
  };

  return (
    <Modal className={modalAdjustments} title={intl.get('edit_report_group')} isOpen={isOpen} onDismiss={onClose}>
      <div
        className={css`
          display: flex;
          flex: 1;
          justify-content: flex-end;
        `}
      >
        <EditReportGroupForm onUpdate={onUpdateReportGroup} reportGroup={reportGroup} />
        <Button size="md" variant="destructive" onClick={setDeleteAlertIsOpen}>
          {intl.get('delete')}
        </Button>
      </div>

      <ReportGroupMemberList reportGroup={reportGroup} />
      <ConfirmModal
        isOpen={deleteAlertIsOpen}
        title={intl.get('delete_report_group')}
        body={intl.get('delete_report_group_question')}
        confirmText={intl.get('delete')}
        icon="exclamation-triangle"
        onConfirm={onConfirmDeleteGroup}
        onDismiss={setDeleteAlertIsOpen}
      />
    </Modal>
  );
};
