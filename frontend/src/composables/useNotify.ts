import { useMessage, useNotification } from 'naive-ui'

export function useNotify() {
  const message = useMessage()
  const notification = useNotification()
  
  return {
    success: (content: string) => message.success(content),
    error: (content: string) => message.error(content),
    warning: (content: string) => message.warning(content),
    info: (content: string) => message.info(content),
    
    notify: {
      success: (title: string, content?: string) => notification.success({ 
        title, 
        content, 
        duration: 3000 
      }),
      error: (title: string, content?: string) => notification.error({ 
        title, 
        content, 
        duration: 5000 
      }),
      info: (title: string, content?: string) => notification.info({ 
        title, 
        content, 
        duration: 3000 
      })
    }
  }
}
