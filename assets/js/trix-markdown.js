// trix-markdown-editor-handler.js
(function() {
  'use strict';

  class TrixMarkdownEditor {
    constructor(editorElement, hiddenInput) {
      this.editor = editorElement;
      this.hiddenInput = hiddenInput;
      this.init();
    }

    init() {
      this.editor.addEventListener('trix-change', () => {
        this.updateMarkdown();
      });

      if (this.hiddenInput.value) {
        this.editor.editor.loadHTML(this.markdownToHTML(this.hiddenInput.value));
      }
    }

    updateMarkdown() {
      const trixDocument = this.editor.editor.getDocument();
      const markdown = this.trixDocumentToMarkdown(trixDocument);
      this.hiddenInput.value = markdown;
    }

    trixDocumentToMarkdown(document) {
      let markdown = '';
      const blocks = document.getBlocks();
      
      for (let i = 0; i < blocks.length; i++) {
        const block = blocks[i];
        const blockMarkdown = this.convertBlock(block);
        
        if (!blockMarkdown.trim()) {
          continue;
        }
        
        markdown += blockMarkdown;
        
        if (i < blocks.length - 1) {
          const nextBlock = blocks[i + 1];
          if (nextBlock && this.convertBlock(nextBlock).trim()) {
            markdown += '\n\n';
          }
        }
      }
      
      return markdown.trim().replace(/\n{3,}/g, '\n\n');
    }

    convertBlock(block) {
      const attributes = block.getAttributes();
      const text = block.getText();
      
      if (!text.trim()) {
        return '';
      }
      
      if (attributes.heading1) {
        return `# ${this.convertInlineFormatting(text, block)}`;
      } else if (attributes.heading2) {
        return `## ${this.convertInlineFormatting(text, block)}`;
      } else if (attributes.heading3) {
        return `### ${this.convertInlineFormatting(text, block)}`;
      } else if (attributes.heading4) {
        return `#### ${this.convertInlineFormatting(text, block)}`;
      } else if (attributes.heading5) {
        return `##### ${this.convertInlineFormatting(text, block)}`;
      } else if (attributes.heading6) {
        return `###### ${this.convertInlineFormatting(text, block)}`;
      } else if (attributes.quote) {
        const content = this.convertInlineFormatting(text, block);
        const lines = content.split('\n');
        return lines.map(line => line.trim() ? `> ${line}` : '>').join('\n');
      } else if (attributes.code) {
        return `\`\`\`\n${text}\n\`\`\``;
      } else if (attributes.bulletList) {
        const content = this.convertInlineFormatting(text, block);
        return `- ${content}`;
      } else if (attributes.numberList) {
        const content = this.convertInlineFormatting(text, block);
        return `1. ${content}`;
      } else {
        const formattedText = this.convertInlineFormatting(text, block);
        return formattedText;
      }
    }

    convertInlineFormatting(text, block) {
      if (!text) return '';
      
      let result = '';
      const pieces = block.getPieces();
      
      for (const piece of pieces) {
        const pieceText = piece.getText();
        const attributes = piece.getAttributes();
        
        if (!pieceText) continue;
        
        let formattedText = pieceText;
        
        if (attributes.href) {
          const linkText = pieceText.trim();
          const href = attributes.href.trim();
          formattedText = `[${linkText}](${href})`;
        } else {
          if (attributes.code) {
            formattedText = `\`${pieceText}\``;
          } else {
            if (attributes.bold && attributes.italic) {
              formattedText = `***${formattedText}***`;
            } else if (attributes.bold) {
              formattedText = `**${formattedText}**`;
            } else if (attributes.italic) {
              formattedText = `*${formattedText}*`;
            }
            
            if (attributes.strike) {
              formattedText = `~~${formattedText}~~`;
            }
          }
        }
        
        result += formattedText;
      }
      
      return result;
    }

  document.addEventListener('DOMContentLoaded', function() {
    const trixEditors = document.querySelectorAll('trix-editor[data-markdown-target]');
    
    trixEditors.forEach(editor => {
      const targetId = editor.getAttribute('data-markdown-target');
      const hiddenInput = document.getElementById(targetId);
      
      if (hiddenInput) {
        new TrixMarkdownEditor(editor, hiddenInput);
        console.log(`TrixMarkdownEditor initialized for editor targeting #${targetId}`);
      } else {
        console.warn(`TrixMarkdownEditor: Hidden input with ID "${targetId}" not found for editor.`);
      }
    });
  });

  window.TrixMarkdownEditor = TrixMarkdownEditor; // Still good for debugging
})();
