#!/usr/bin/env python3
"""
Python 实现：启动10万个线程循环打印
对比 Go 的 goroutine 实现
"""

import threading
import time
from concurrent.futures import ThreadPoolExecutor, as_completed
import asyncio


def method1_threading():
    """
    方法1: 使用 threading.Thread 直接创建10万个线程
    注意：Python 线程受 GIL 限制，且创建10万个真实线程可能会耗尽系统资源
    """
    print("方法1: 使用 threading.Thread (不推荐，可能会失败)")
    
    threads = []
    lock = threading.Lock()  # 用于同步打印，避免输出混乱
    
    def worker(i):
        # 使用 lock 避免打印混乱
        with lock:
            print(i, end=' ')
    
    start_time = time.time()
    
    try:
        # 创建并启动10万个线程
        for i in range(100000):
            t = threading.Thread(target=worker, args=(i,))
            threads.append(t)
            t.start()
        
        # 等待所有线程完成
        for t in threads:
            t.join()
        
        print(f"\n方法1完成，耗时: {time.time() - start_time:.2f}秒")
    except Exception as e:
        print(f"\n方法1失败: {e}")


def method2_threadpool():
    """
    方法2: 使用 ThreadPoolExecutor 线程池
    推荐方式：使用线程池管理线程，避免创建过多线程
    """
    print("\n方法2: 使用 ThreadPoolExecutor 线程池")
    
    def worker(i):
        print(i, end=' ')
        return i
    
    start_time = time.time()
    
    # 使用线程池，设置合理的最大线程数
    # 通常设置为 CPU 核心数的 2-4 倍
    max_workers = 1000  # 可以根据系统调整
    
    with ThreadPoolExecutor(max_workers=max_workers) as executor:
        # 提交10万个任务
        futures = [executor.submit(worker, i) for i in range(100000)]
        
        # 等待所有任务完成
        for future in as_completed(futures):
            pass  # 任务已经在 worker 中打印
    
    print(f"\n方法2完成，耗时: {time.time() - start_time:.2f}秒")


async def method3_asyncio():
    """
    方法3: 使用 asyncio 协程
    最接近 Go goroutine 的实现方式，轻量级并发
    """
    print("\n方法3: 使用 asyncio 协程 (最接近 Go goroutine)")
    
    async def worker(i):
        print(i, end=' ')
        # 模拟异步操作
        await asyncio.sleep(0)
    
    start_time = time.time()
    
    # 创建10万个协程任务
    tasks = [worker(i) for i in range(100000)]
    
    # 并发执行所有任务
    await asyncio.gather(*tasks)
    
    print(f"\n方法3完成，耗时: {time.time() - start_time:.2f}秒")


def method4_asyncio_semaphore():
    """
    方法4: 使用 asyncio + Semaphore 控制并发数
    推荐方式：既高效又不会耗尽资源
    """
    print("\n方法4: 使用 asyncio + Semaphore (推荐)")
    
    async def run_with_semaphore():
        # 使用信号量控制同时运行的协程数量
        semaphore = asyncio.Semaphore(1000)
        
        async def worker(i):
            async with semaphore:
                print(i, end=' ')
                await asyncio.sleep(0)
        
        start_time = time.time()
        
        # 创建10万个协程任务
        tasks = [worker(i) for i in range(100000)]
        
        # 并发执行所有任务
        await asyncio.gather(*tasks)
        
        print(f"\n方法4完成，耗时: {time.time() - start_time:.2f}秒")
    
    asyncio.run(run_with_semaphore())


def main():
    print("=" * 60)
    print("Python vs Go: 10万个并发任务打印数字")
    print("=" * 60)
    
    choice = input("""
选择运行方式:
1. threading.Thread (不推荐，可能失败)
2. ThreadPoolExecutor (推荐用于 I/O 密集型)
3. asyncio 协程 (推荐，最接近 Go goroutine)
4. asyncio + Semaphore (推荐，资源控制)
5. 仅运行 asyncio 协程版本

请输入选项 (1-5, 默认5): """).strip() or "5"
    
    if choice == "1":
        method1_threading()
    elif choice == "2":
        method2_threadpool()
    elif choice == "3":
        asyncio.run(method3_asyncio())
    elif choice == "4":
        method4_asyncio_semaphore()
    elif choice == "5":
        print("\n运行 asyncio 协程版本 (最接近 Go 实现):")
        asyncio.run(method3_asyncio())
    else:
        print("无效选项")
    
    print("\n" + "=" * 60)
    print("对比说明:")
    print("- Go goroutine: 轻量级，由 Go runtime 调度，可以轻松创建百万级")
    print("- Python Thread: 真实线程，受 GIL 和系统资源限制")
    print("- Python asyncio: 协程，轻量级，最接近 goroutine 的实现")
    print("=" * 60)


if __name__ == "__main__":
    main()

