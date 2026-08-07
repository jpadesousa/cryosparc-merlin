#!/usr/bin/env bash
# cryoSPARC cluster submission script template for SLURM
# Lane: gpu-hourly

{# This template uses jinja2 syntax. #}

# Available variables:
#  script_path_abs={{ script_path_abs }}
#      - the absolute path to the generated submission script
#  run_cmd={{ run_cmd }}
#      - the complete command-line string to run the job
#  num_cpu={{ num_cpu }}
#      - the number of CPUs needed
#  num_gpu={{ num_gpu }}
#      - the number of GPUs needed. Note: the code will use this many GPUs
#        starting from dev id 0. The cluster scheduler or this script have the
#        responsibility of setting CUDA_VISIBLE_DEVICES so that the job code
#        ends up using the correct cluster-allocated GPUs.
#  ram_gb={{ ram_gb }}
#      - the amount of RAM needed in GB
#  job_dir_abs={{ job_dir_abs }}
#      - absolute path to the job directory
#  project_dir_abs={{ project_dir_abs }}
#      - absolute path to the project dir
#  job_log_path_abs={{ job_log_path_abs }}
#      - absolute path to the log file for the job
#  worker_bin_path={{ worker_bin_path }}
#      - absolute path to the cryosparc worker command
#  run_args={{ run_args }}
#      - arguments to be passed to cryosparcw run
#  project_uid={{ project_uid }}
#      - uid of the project
#  job_uid={{ job_uid }}
#      - uid of the job
#  job_creator={{ job_creator }}
#      - name of the user that created the job (may contain spaces)
#  cryosparc_username={{ cryosparc_username }}
#      - cryosparc username of the user that created the job (usually an email)

#SBATCH --job-name=cryosparc_{{ project_uid }}_{{ job_uid }}
#SBATCH --output={{ job_log_path_abs }}.out
#SBATCH --error={{ job_log_path_abs }}.err
#SBATCH --ntasks=1
#SBATCH --mem={{ (ram_gb*1000*2)|int }}M
#SBATCH --time=01:00:00
#SBATCH --partition=a100-hourly
#SBATCH --cluster=gmerlin7
#SBATCH --gres=gpu:{{ num_gpu }}
#SBATCH --cpus-per-task={{ (num_gpu*16) }}
#SBATCH --hint=nomultithread

{%- if num_gpu == 0 %}
# Use CPU cluster
echo "Error: No GPU requested. Use a CPU lane instead." >&2
exit 1
{%- else %}

# Print hostname, for debugging
echo "Job Id: $SLURM_JOBID"
echo "Host:   $SLURM_NODELIST"

module purge

srun {{ run_cmd }}

EXIT_CODE=$?
echo "Exit code: $EXIT_CODE"
exit $?
{%- endif %}