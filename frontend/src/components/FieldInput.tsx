import React, { FC } from 'react';

import { FieldValidationMessage, InlineFormLabel, Input } from '@grafana/ui';

type Props = {
  label: string;
  defaultValue: string | number;
  placeholder?: string;
  inputName: string;
  invalid: boolean;
  errorMessage: string;
  tooltip: string;
  register: any;
  type?: string;
};

export const FieldInput = React.forwardRef<HTMLInputElement, Props>(
  ({ tooltip, type = 'text', label, register, defaultValue, placeholder, inputName, invalid, errorMessage }, ref) => {
    return (
      <div className="gf-form">
        <InlineFormLabel className="width-14" tooltip={tooltip}>
          {label}
        </InlineFormLabel>
        <div style={{ display: 'flex', flex: 1, flexDirection: 'column' }}>
          <Input
            type={type}
            defaultValue={defaultValue}
            placeholder={placeholder}
            name={inputName}
            ref={register()}
            css=""
          />

          {invalid && (
            <div>
              <FieldValidationMessage>{errorMessage}</FieldValidationMessage>
            </div>
          )}
        </div>
      </div>
    );
  }
);
