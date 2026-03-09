import { Body, Controller, Get, Param, Post } from '@nestjs/common';
import { OrdersService } from './orders.service';
import { CreateOrderDto } from './dto/create-order.dto';
import type { OrderResponse } from './order.response';
import { toOrderResponse } from './order.response';

@Controller('orders')
export class OrdersController {
  constructor(private readonly ordersService: OrdersService) {}

  @Post()
  create(@Body() dto: CreateOrderDto): OrderResponse {
    return toOrderResponse(this.ordersService.create(dto));
  }

  @Get(':id')
  findOne(@Param('id') id: string): OrderResponse {
    return toOrderResponse(this.ordersService.findById(id));
  }

  @Get()
  findAll(): OrderResponse[] {
    return this.ordersService.findAll().map(toOrderResponse);
  }
}
