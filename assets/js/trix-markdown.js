// Trix Editor with Markdown Conversion
(function() {
  'use strict';

  class TrixMarkdownEditor {
    constructor(editorElement, hiddenInput) {
      this.editor = editorElement;
      this.hiddenInput = hiddenInput;
      this.init();
    }

    init() {
      // Listen for content changes
      this.editor.addEventListener('trix-change', () => {
        this.updateMarkdown();
      });

      // Initialize with existing content if any
      if (this.hiddenInput.value) {
        this.editor.editor.loadHTML(this.markdownToHTML(this.hiddenInput.value));
      }
    }

    updateMarkdown() {
      const html = this.editor.innerHTML;
      const markdown = this.htmlToMarkdown(html);
      this.hiddenInput.value = markdown;
    }

    htmlToMarkdown(html) {
      // Create a temporary div to parse HTML
      const temp = document.createElement('div');
      temp.innerHTML = html;

      // Convert HTML to Markdown
      return this.processNode(temp);
    }

    processNode(node) {
      let markdown = '';

      for (const child of node.childNodes) {
        if (child.nodeType === Node.TEXT_NODE) {
          markdown += child.textContent;
        } else if (child.nodeType === Node.ELEMENT_NODE) {
          markdown += this.convertElement(child);
        }
      }

      return markdown;
    }

    convertElement(element) {
      const tagName = element.tagName.toLowerCase();
      const content = this.processNode(element);

      switch (tagName) {
        case 'h1':
          return `# ${content}\n\n`;
        case 'h2':
          return `## ${content}\n\n`;
        case 'h3':
          return `### ${content}\n\n`;
        case 'h4':
          return `#### ${content}\n\n`;
        case 'h5':
          return `##### ${content}\n\n`;
        case 'h6':
          return `###### ${content}\n\n`;
        case 'p':
          return `${content}\n\n`;
        case 'strong':
        case 'b':
          return `**${content}**`;
        case 'em':
        case 'i':
          return `*${content}*`;
        case 'u':
          return `<u>${content}</u>`;
        case 'del':
        case 's':
          return `~~${content}~~`;
        case 'a':
          const href = element.getAttribute('href');
          return href ? `[${content}](${href})` : content;
        case 'ul':
          return this.convertList(element, false);
        case 'ol':
          return this.convertList(element, true);
        case 'li':
          return content;
        case 'blockquote':
          return content.split('\n').map(line => line.trim() ? `> ${line}` : '>').join('\n') + '\n\n';
        case 'pre':
          return `\`\`\`\n${content}\n\`\`\`\n\n`;
        case 'code':
          return `\`${content}\``;
        case 'br':
          return '\n';
        case 'div':
          // Handle Trix block breaks
          if (element.innerHTML === '<br>') {
            return '\n';
          }
          return content;
        default:
          return content;
      }
    }

    convertList(listElement, isOrdered) {
      let markdown = '';
      let counter = 1;

      for (const li of listElement.children) {
        if (li.tagName.toLowerCase() === 'li') {
          const content = this.processNode(li).trim();
          const prefix = isOrdered ? `${counter}. ` : '- ';
          markdown += `${prefix}${content}\n`;
          if (isOrdered) counter++;
        }
      }

      return markdown + '\n';
    }

    markdownToHTML(markdown) {
      // Basic markdown to HTML conversion for initial loading
      let html = markdown;

      // Headers
      html = html.replace(/^### (.*$)/gim, '<h3>$1</h3>');
      html = html.replace(/^## (.*$)/gim, '<h2>$1</h2>');
      html = html.replace(/^# (.*$)/gim, '<h1>$1</h1>');

      // Bold
      html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');

      // Italic
      html = html.replace(/\*(.*?)\*/g, '<em>$1</em>');

      // Links
      html = html.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2">$1</a>');

      // Line breaks
      html = html.replace(/\n\n/g, '</p><p>');
      html = html.replace(/\n/g, '<br>');

      // Wrap in paragraphs
      if (html && !html.startsWith('<')) {
        html = '<p>' + html + '</p>';
      }

      return html;
    }
  }

  // Initialize Trix editors when DOM is loaded
  document.addEventListener('DOMContentLoaded', function() {
    const trixEditors = document.querySelectorAll('trix-editor[data-markdown-target]');
    
    trixEditors.forEach(editor => {
      const targetId = editor.getAttribute('data-markdown-target');
      const hiddenInput = document.getElementById(targetId);
      
      if (hiddenInput) {
        new TrixMarkdownEditor(editor, hiddenInput);
      }
    });
  });

  // Make available globally
  window.TrixMarkdownEditor = TrixMarkdownEditor;
})();
