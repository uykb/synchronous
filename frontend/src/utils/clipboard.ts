import { ref } from 'vue'

export function useClipboard() {
  const copied = ref(false)
  
  const copyToClipboard = async (text: string): Promise<boolean> => {
    if (!text) return false
    
    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text)
      } else {
        // Fallback for non-secure contexts
        const textArea = document.createElement('textarea')
        textArea.value = text
        textArea.style.cssText = 'position:fixed;left:-9999px;top:0'
        document.body.appendChild(textArea)
        textArea.focus()
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
      }
      
      copied.value = true
      setTimeout(() => { copied.value = false }, 2000)
      return true
    } catch (err) {
      console.error('Copy failed:', err)
      return false
    }
  }
  
  return { copied, copyToClipboard }
}
