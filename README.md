### Project structure
based on pattern: <br>
https://github.com/evrone/go-clean-template <br>
https://evrone.ru/Go-clean-template <br>
https://www.youtube.com/watch?v=V6lQG6d5LgU


`/cmd/app/main.go`
Чтение конфигурационных файлов, переменных окружения и старт м/сервиса запуском процедуры в файле /internal/app/initapp.go

`/internal/app/initapp.go`
Создание (конструктор New) и запуск основных компонентов м/сервиса. Передача в них зависимостей (Dependency injection)

## Основные компоненты
- БД postgres (реализация в pkg/pg/pg.go). Передаем в нее зависимость - структуру cfg
- Структура taskUseCase. Основа обработки бизнес логики. Инкапсулирует через поле dbRepo интерфейс работы с БД. Сама структура имплементирует интерфейс (type TaskHandler interface) для вызова ее из методов роутера.
- http server - стандартный сервер (реализация в pkg/httpserver/httpserver.go) и стандартный роутер. Создаем роутер и передаём его вместе со структурой cfg в реализацию сервера

### Сущности
`/internal/entity/`
Здесь располагаются сущности бизнес-модели

### Бизнес логика
`/internal/usecase/`
Бизнес логика. Не зависит от реализации репозитория.
InUseCase - обрабатывает входящие запросы (общается с фронтом - роутером через интерфейс). Реализует методы TaskHandlerInterface, которые "дергаются" из роутера.
В InUseCase встроены другие юзкейсы - TaskUseCase и UserUseCase, которые общаются с back-ом через интерфейсы TaskDBRepoInterface и UserDBRepoInterface соответственно.

`/internal/usecase/interfaces.go`
- интерфейсы работы с БД (type TaskDBRepo interface), которые реализуются в файлах репозитория, а "дергаются" в файлах этого пакета (usecase) 
- реализация интерфейса репозитория для БД postgresql в файле /internal/repository/task_pg_repo.go
                                    для простой записи в файл - /internal/repository/task_mock_repo.go

### Репозиторий
`/internal/repository/`
Здесь располагаются файлы работы с БД и kafka.
- В случае с БД в репозиторий получаем зависимости - БД postgres и модели - entities. Непосредственные CRUDL-операции с БД. Располагающиеся здесь методы реализуют интерфейсы работы с БД, объявленные в usecase package in /internal/usecase/interfaces.go  (type TaskDBRepo interface)
- В случае с kafka аналогично

### Входящие запросы
`/internal/controller/grpc`
Методы клиента grpc, который "дергает" сервис auth и проверяет валидность токена user-а

`/internal/controller/http/v1/`
Router<br>
Вначале проверка на валидность токена через gRPC.<br>
После этого обработка entrypoint-ов - вызов методов пакета usecase через интерфейс (type TaskHandler interface).

`/internal/controller/http/v1/interfaces.go`
Интерфейсы работы с юзкейсами (type TaskHandler interface), которые реализуются в файлах пакета usecase, а "дергаются" в методах http controller-а (/internal/controller/http/v1/task_rout.go)

### GRPC
`/api/grpc/proto` - здесь контракт
`/api/grpc/gen` - здесь сгенерированные файлы
### Логирование:
`/pkg/logging/applogging.go`
Логгер на базе стандартного логера. Для каждого уровня логирования создан свой логгер. Цель - возможность записывать логи каждого уровня в разные файлы. Наверняка это позволяют делать и другие логгеры, но пока реализовано на стандартном. <br>
Вывод будет осуществляться в три файла - Info (level Info), Error (levels Warn, Error, Fatal) и Debug (level Debug) - это задается в конструкторе NewTaskLogger.
Куда будут выводиться логи задается в config-файле (секция Log) задается путь к файлу. По умолчанию Info и Debug levels в StdOut, а остальные в StdErr.


## Request bodys

### Create
`{
	"descr": "descr",
	"body": "body",
	"approvers": [{
			"email": "appr1@gmail.com"
		},
		{
			"email": "appr2@gmail.com"
		},
		{
			"email": "appr3@gmail.com"
		},
		{
			"email": "appr4@gmail.com"
		}
	]
}`

### Update (PUT)  // id in path
`{
	"descr": "descr",
	"body": "body",
	"approvers": [{
			"email": "appr1@mail.ru"
		},
		{
			"email": "appr2@mail.ru"
		},
		{
			"email": "appr3@mail.ru"
		},
		{
			"email": "appr4@mail.ru"
		}
	]
}`
