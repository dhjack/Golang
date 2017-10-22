#include <stdlib.h>
#include <stdio.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/shm.h>
#include <errno.h>

void* getShareMemory(int shm_key, int size) {

    void* shm_p = NULL;

    int shmid = shmget(shm_key, size, 0666 | IPC_CREAT | IPC_EXCL);
    if (shmid == -1) {
        if (EEXIST == errno) {
            printf("sharemem mem is exist!\n");
            shmid = shmget(shm_key, size, 0666);
            if (shmid == -1) {
                perror("sharemem get failed and call exit");
                exit(EXIT_FAILURE);
            }
        } else {
            perror("sharemem get failed with flag IPC_CREAT IPC_EXCL and call exit!");
            exit(EXIT_FAILURE);
        }
    }

    printf("sharemem shmid:%d\n", shmid);

    shm_p = shmat(shmid, NULL, 0);

    if (!shm_p) {
        perror("shremem bind failed and call exit");
        exit(EXIT_FAILURE);
    }

    printf("sharemem shmid:%p\n", shm_p);
    return shm_p;
}

