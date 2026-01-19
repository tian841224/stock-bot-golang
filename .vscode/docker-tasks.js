#!/usr/bin/env node

const { execSync } = require('child_process');
const path = require('path');

// 取得命令列參數
const args = process.argv.slice(2);
const command = args[0];

// 顏色輸出函數
const colors = {
    reset: '\x1b[0m',
    bright: '\x1b[1m',
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    magenta: '\x1b[35m',
    cyan: '\x1b[36m'
};

function log(message, color = 'reset') {
    console.log(`${colors[color]}${message}${colors.reset}`);
}

function executeCommand(command, description) {
    try {
        log(`\n${colors.cyan}執行: ${description}${colors.reset}`);
        log(`${colors.yellow}命令: ${command}${colors.reset}\n`);
        
        const output = execSync(command, { 
            stdio: 'inherit',
            cwd: process.cwd()
        });
        
        log(`\n${colors.green}✓ ${description} 完成${colors.reset}\n`);
        return true;
    } catch (error) {
        log(`\n${colors.red}✗ ${description} 失敗: ${error.message}${colors.reset}\n`);
        return false;
    }
}

function checkDockerRunning() {
    try {
        execSync('docker --version', { stdio: 'pipe' });
        execSync('docker-compose --version', { stdio: 'pipe' });
        return true;
    } catch (error) {
        log(`${colors.red}錯誤: Docker 或 Docker Compose 未安裝或未執行${colors.reset}`);
        return false;
    }
}

function checkEnvFile() {
    const fs = require('fs');
    const envPath = path.join(process.cwd(), '.env');
    
    if (!fs.existsSync(envPath)) {
        log(`${colors.yellow}警告: .env 檔案不存在，請確認環境變數設定${colors.reset}`);
        return false;
    }
    return true;
}

// 主要命令處理
switch (command) {
    case 'start-all':
        log(`${colors.bright}${colors.blue}🐳 啟動所有 Docker 服務${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        checkEnvFile();
        
        executeCommand(
            'docker-compose up -d --build',
            '啟動所有服務 (PostgreSQL + Stock Bot + Sync Service)'
        );
        
        log(`${colors.green}所有服務已啟動！${colors.reset}`);
        log(`${colors.cyan}• PostgreSQL: localhost:5432${colors.reset}`);
        log(`${colors.cyan}• Stock Bot: localhost:8080${colors.reset}`);
        log(`${colors.cyan}• 查看日誌: docker-compose logs -f${colors.reset}`);
        break;

    case 'start-bot':
        log(`${colors.bright}${colors.blue}🤖 啟動 Stock Bot 服務${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        checkEnvFile();
        
        executeCommand(
            'docker-compose up -d postgres',
            '啟動 PostgreSQL 資料庫'
        );
        
        executeCommand(
            'docker-compose up -d stock-bot',
            '啟動 Stock Bot 應用程式'
        );
        
        log(`${colors.green}Stock Bot 已啟動！${colors.reset}`);
        log(`${colors.cyan}• 服務網址: localhost:8080${colors.reset}`);
        break;

    case 'start-debug':
        log(`${colors.bright}${colors.blue}🐛 啟動除錯模式${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        checkEnvFile();
        
        executeCommand(
            'docker-compose -f docker-compose_debug.yml up -d --build',
            '啟動除錯模式服務'
        );
        
        log(`${colors.green}除錯模式已啟動！${colors.reset}`);
        log(`${colors.cyan}• 使用 docker-compose_debug.yml 設定${colors.reset}`);
        log(`${colors.cyan}• 查看日誌: docker-compose -f docker-compose_debug.yml logs -f${colors.reset}`);
        break;

    case 'stop-all':
        log(`${colors.bright}${colors.red}🛑 停止所有 Docker 服務${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        executeCommand(
            'docker-compose down',
            '停止所有服務'
        );
        
        log(`${colors.green}所有服務已停止！${colors.reset}`);
        break;

    case 'logs':
        log(`${colors.bright}${colors.blue}📋 查看服務日誌${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        const service = args[1] || '';
        const logCommand = service ? 
            `docker-compose logs -f ${service}` : 
            'docker-compose logs -f';
            
        executeCommand(logCommand, `查看 ${service || '所有服務'} 日誌`);
        break;

    case 'status':
        log(`${colors.bright}${colors.blue}📊 服務狀態${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        executeCommand(
            'docker-compose ps',
            '查看服務狀態'
        );
        break;

    case 'clean':
        log(`${colors.bright}${colors.yellow}🧹 清理 Docker 資源${colors.reset}`);
        
        if (!checkDockerRunning()) {
            process.exit(1);
        }
        
        executeCommand(
            'docker-compose down -v --remove-orphans',
            '停止並移除所有容器和卷'
        );
        
        executeCommand(
            'docker system prune -f',
            '清理未使用的 Docker 資源'
        );
        
        log(`${colors.green}清理完成！${colors.reset}`);
        break;

    default:
        log(`${colors.bright}Docker 任務助手${colors.reset}`);
        log(`${colors.cyan}可用命令:${colors.reset}`);
        log(`  ${colors.green}start-all${colors.reset}    - 啟動所有服務`);
        log(`  ${colors.green}start-bot${colors.reset}     - 只啟動 Bot 服務`);
        log(`  ${colors.green}start-debug${colors.reset}   - 啟動除錯模式`);
        log(`  ${colors.green}stop-all${colors.reset}      - 停止所有服務`);
        log(`  ${colors.green}logs${colors.reset}          - 查看日誌 (可指定服務名稱)`);
        log(`  ${colors.green}status${colors.reset}        - 查看服務狀態`);
        log(`  ${colors.green}clean${colors.reset}          - 清理 Docker 資源`);
        log(`\n${colors.yellow}範例:${colors.reset}`);
        log(`  node docker-tasks.js start-all`);
        log(`  node docker-tasks.js logs stock-bot`);
        break;
}
