# Opinionated cosmos cron module

This project has two modules:
- cosmos-quartz (default)
- cron

## Cron module

This module is a simple cron module that can be used to schedule jobs.  
It uses scheduler integrated with Keeper and Module and hence it is  
available in all modules inside application.

## Why it is opinionated?

- By requirement, there's no need to implement queries and messages for task scheduling
- Hence, the all request and response are handled by other modules but this implementation 
examples falls out of scope of this test assignment.
- The State and Keeper, they do not hold the state of the task.  

Tasks that are distributed among multiple nodes, a single-node implementation like 
the one shown may not be sufficient.   
In this case, you may want to consider using implementation of a predefined function in the Keeper.

This is how persisting of the task might work:
- Implement function in Keeper that update the State
- Create a task struct type that would keep the cron task, function name and parameters can be serialized.

But aforementioned implementation is not specified in requirements.  
Moreover, it gives freedom to schedule just any function in the application.
But this implementation is not serializable and hence it is not possible to store it in the State.

## Tests
I'm sorry, I must implement tests, but I'm very constrained in time.