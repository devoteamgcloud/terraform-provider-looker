# Notes

## Workflows


```mermaid
flowchart LR;
  subgraph User
    direction LR
    email[Email] --> user(User)
    fname[Firstname] --> user
    lname[Lastame] --> user
  end
  User --->|creates| userid{{User ID}}
  
  subgraph Group
    direction LR
    groupName[Group Name] --> group(Group)
  end
  Group --->|creates| groupid{{Group ID}}
  
  subgraph GroupUsersBinding
    direction LR
    egm[Group Member]
  end
  userid --> GroupUsersBinding
  groupid --> GroupUsersBinding
  
  
  subgraph GroupGroupBinding
    direction LR
    ggm[Group Member]
  end
  groupid --> GroupGroupBinding
  GroupGroupBinding --> groupid
  %%ggm --> groupid
  
  
  
  subgraph Role
    direction LR
    roleName[Role Name] --> role(Role)
  end
  Role -->|creates| roleid{{Role ID}}
  
  subgraph RoleGroupsBinding
    direction LR
  end
  
  %% ModelSet {name, models[]string} > Creates model_set_id
  subgraph ModelSet
    direction LR
    name[ModelSet Name] --> mset
    subgraph Models
      model1[Model]
      model2[Model]
    end
    Models --> mset(Model Set)
  end
  
  ModelSet -->|creates| modelsetid{{Model Set ID}}
  
  subgraph PermissionSet
    direction LR
    pset[Permission Set]
  end
  
  PermissionSet -->|has static| permsetid{{Permission Set ID}}
  
  modelsetid --> Role
  permsetid --> Role
  
  
  subgraph DBConnection
    dbconname[DB Connection name]
  end
  
  
  subgraph Project
    projectname[Project Name]
  end
  
  Project --> |has project name| modelproject  
  
  subgraph LookMlModel
    direction LR
    lookmlmodelname[LookML Model Name]
    modelproject[Project Name]
    subgraph AllowedDBConnNames
    direction LR
      dbconname1[DB Connection Name 1]
      dbconname2[DB Connection Name 2]
    end
  end
  
  LookMlModel --> |has model name| model1
  DBConnection --> |???| dbconname1
  
  
```

```mermaid
%% Example
flowchart TB;

  User -->|creates| editoruserid{{Editor User ID}}
  User -->|creates| adminuserid{{Admin User ID}}
  
  Group -->|creates| editorgroupid{{Editor Group ID}}
  Group -->|creates| admingroupid{{Admin Group ID}}
  
  subgraph GroupUsersBindingEditorGroupBindingUser
    direction LR
    egm[Group Member]
  end
  editoruserid --> EditorGroupBindingUser
  editorgroupid --> EditorGroupBindingUser
  
  subgraph AdminGroupBindingUser
    direction LR
    agm[Group Member]
  end
  adminuserid --> AdminGroupBindingUser
  admingroupid --> AdminGroupBindingUser

```

ddd
xxx

```mermaid
  subgraph GroupBindingRole
    direction LR
    groupmemberrole[Role Member]
  end
  roleid --> GroupBindingRole
  groupid --> GroupBindingRole
  permset(Permission Set) --->|has| permsetid{{Permission Set ID}}

flowchart TB;
  subgraph User
    direction LR
    email[Email] --> user
    fname[Firstname] --> user
    lname[Lastame] --> user
  end
  User -->x
```


```mermaid
flowchart TB
    c1-->a2
    subgraph one
    a1-->a2
    end
    subgraph two
    b1-->b2
    end
    subgraph three
    c1-->c2
    end
```


```mermaid
flowchart LR
  subgraph TOP
    direction TB
    subgraph B1
        direction RL
        i1 -->f1
    end
    subgraph B2
        direction TB
        i2
        f2
    end
  end
  A --> TOP --> B
  B1 --> B2
```


```mermaid
flowchart LR
  subgraph TOP
    direction TB
    subgraph B1
        direction RL
        i1 -->f1
    end
    subgraph B2
        direction BT
        i2 -->f2
    end
  end
  A --> TOP --> B
  B1 --> B2
```