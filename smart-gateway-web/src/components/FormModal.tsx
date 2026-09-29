import { useEffect } from 'react'
import { Form, Input, InputNumber, Modal, Select } from 'antd'
import type { ReactNode } from 'react'
import type { Option } from '@/constants/options'

export type FieldType = 'text' | 'password' | 'number' | 'select'

export interface FieldConfig {
  name: string
  label: string
  type: FieldType
  required?: boolean
  options?: Option[]
  placeholder?: string
  min?: number
  max?: number
  addonAfter?: string
  disabled?: boolean
}

interface FormModalProps {
  open: boolean
  title: string
  fields: FieldConfig[]
  initialValues?: Record<string, unknown>
  confirmLoading?: boolean
  onCancel: () => void
  onSubmit: (values: Record<string, unknown>) => void
}

export default function FormModal({
  open,
  title,
  fields,
  initialValues,
  confirmLoading,
  onCancel,
  onSubmit,
}: FormModalProps) {
  const [form] = Form.useForm()

  useEffect(() => {
    if (open) {
      form.resetFields()
      if (initialValues) form.setFieldsValue(initialValues)
    }
  }, [open, initialValues, form])

  const renderField = (f: FieldConfig): ReactNode => {
    switch (f.type) {
      case 'number':
        return (
          <InputNumber
            style={{ width: '100%' }}
            min={f.min}
            max={f.max}
            placeholder={f.placeholder}
            addonAfter={f.addonAfter}
          />
        )
      case 'select':
        return <Select options={f.options} placeholder={f.placeholder} allowClear />
      case 'password':
        return <Input.Password placeholder={f.placeholder} autoComplete="new-password" />
      default:
        return <Input placeholder={f.placeholder} disabled={f.disabled} />
    }
  }

  const handleOk = async () => {
    const values = await form.validateFields()
    onSubmit(values)
  }

  return (
    <Modal open={open} title={title} onCancel={onCancel} onOk={handleOk} confirmLoading={confirmLoading}>
      <Form form={form} layout="vertical" style={{ marginTop: 8 }}>
        {fields.map((f) => (
          <Form.Item
            key={f.name}
            name={f.name}
            label={f.label}
            rules={f.required ? [{ required: true, message: `请输入${f.label}` }] : undefined}
          >
            {renderField(f)}
          </Form.Item>
        ))}
      </Form>
    </Modal>
  )
}
